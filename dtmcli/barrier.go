package dtmcli

import (
	"database/sql"
	"fmt"
	"net/url"
)

// BusiFunc type for busi func
type BusiFunc func(db DB) error

// BranchBarrier every branch info
type BranchBarrier struct {
	TransType  string
	Gid        string
	BranchID   string
	BranchType string
	BarrierID  int
}

func (bb *BranchBarrier) String() string {
	return fmt.Sprintf("transInfo: %s %s %s %s", bb.TransType, bb.Gid, bb.BranchID, bb.BranchType)
}

// BarrierFromQuery construct transaction info from request
func BarrierFromQuery(qs url.Values) (*BranchBarrier, error) {
	return BarrierFrom(qs.Get("trans_type"), qs.Get("gid"), qs.Get("branch_id"), qs.Get("branch_type"))
}

// BarrierFrom construct transaction info from request
func BarrierFrom(transType, gid, branchID, branchType string) (*BranchBarrier, error) {
	ti := &BranchBarrier{
		TransType:  transType,
		Gid:        gid,
		BranchID:   branchID,
		BranchType: branchType,
	}
	if ti.TransType == "" || ti.Gid == "" || ti.BranchID == "" || ti.BranchType == "" {
		return nil, fmt.Errorf("invlid trans info: %v", ti)
	}
	return ti, nil
}

func insertBarrier(tx Tx, transType string, gid string, branchID string, branchType string, barrierID string, reason string) (int64, error) {
	// 忽略正常语义的动作时要补插入的记录
	if branchType == "" {
		return 0, nil
	}
	return DBExec(tx, "insert ignore into dtm_barrier.barrier(trans_type, gid, branch_id, branch_type, barrier_id, reason) values(?,?,?,?,?,?)", transType, gid, branchID, branchType, barrierID, reason)
}

// Call 子事务屏障，详细介绍见 https://zhuanlan.zhihu.com/p/388444465
// tx: 本地数据库的事务对象，允许子事务屏障进行事务操作
// bisiCall: 业务函数，仅在必要时被调用
func (bb *BranchBarrier) Call(tx Tx, busiCall BusiFunc) (rerr error) {
	bb.BarrierID = bb.BarrierID + 1
	bid := fmt.Sprintf("%02d", bb.BarrierID)
	if rerr != nil {
		return
	}
	defer func() {
		// Logf("barrier call error is %v", rerr)
		if x := recover(); x != nil {
			tx.Rollback()
			panic(x)
		} else if rerr != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	ti := bb
	originType := map[string]string{
		BranchCancel:     BranchTry,    // 用于TCC模式
		BranchCompensate: BranchAction, // 用于SAGA模式
	}[ti.BranchType]
	// 如果是正常语义的动作，originType=""（空字符串），不会插入动作记录（所以affected=0）
	// 如果是撤销语义的动作，originType=“正常语义的动作”，补插入一条正常语义动作的记录
	originAffected, _ := insertBarrier(tx, ti.TransType, ti.Gid, ti.BranchID, originType, bid, ti.BranchType)
	// 如果是正常语义的动作，直接插入
	// 如果是撤销语义的动作，插入撤销语义动作对应的记录
	currentAffected, rerr := insertBarrier(tx, ti.TransType, ti.Gid, ti.BranchID, ti.BranchType, bid, ti.BranchType)
	Logf("originAffected: %d currentAffected: %d", originAffected, currentAffected)
	// 当前是撤销语义的动作，但是正常语义的动作未执行（补插入成功），说明是UNDO先于DO，属于空补偿，直接返回成功
	if (ti.BranchType == BranchCancel || ti.BranchType == BranchCompensate) && originAffected > 0 { // 这个是空补偿，返回成功
		return
	} else if currentAffected == 0 { // 插入不成功
		var result sql.NullString
		err := DBQueryRow(tx, "select 1 from dtm_barrier.barrier where trans_type=? and gid=? and branch_id=? and branch_type=? and barrier_id=? and reason=?",
			ti.TransType, ti.Gid, ti.BranchID, ti.BranchType, bid, ti.BranchType).Scan(&result)
		// 当前是正常语义的动作，但是已经有另一个原因（撤销语义）插入的数据，说明DO晚于UNDO，属于悬挂请求，返回失败
		if err == sql.ErrNoRows { // 不是当前分支插入的，那么是cancel插入的，因此是悬挂操作，返回失败，AP收到这个返回，会尽快回滚
			rerr = ErrFailure
			return
		}
		// 已经有完全相同的动作（gid，bid，actionid，reason）执行过，根据幂等性，直接返回
		rerr = err //幂等和空补偿，直接返回
		return
	}
	// 通过barrier，执行动作逻辑
	rerr = busiCall(tx)
	return
}
