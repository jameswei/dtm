create database if not exists dtm_barrier /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

drop table if exists dtm_barrier.barrier;
create table if not exists dtm_barrier.barrier(
  id int(11) PRIMARY KEY AUTO_INCREMENT,
  trans_type varchar(45) default '' COMMENT '事务类型: saga | xa | tcc | msg',
  gid varchar(128) default'' COMMENT '事务全局id',
  branch_id varchar(128) default '' COMMENT '事务分支名称',
  branch_type varchar(45) default '' COMMENT '事务分支类型 saga_action | saga_compensate | xa',
  barrier_id varchar(45) default '',
  reason varchar(45) default '' comment 'the branch type who insert this record',
  create_time datetime DEFAULT now(),
  update_time datetime DEFAULT now(),
  key(create_time),
  key(update_time),
  UNIQUE key(gid, branch_id, branch_type, barrier_id)
);
