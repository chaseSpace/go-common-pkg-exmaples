drop table test_idxquery_while_longtx_running;
create table test_idxquery_while_longtx_running
(
    id           int primary key auto_increment,
    c            varchar(32) unique,
    nonce_no_idx varchar(32)
);

drop procedure if exists test_idxquery_while_longtx_running_1;
create procedure test_idxquery_while_longtx_running_1()
begin
    set @count = 100000;
#     start transaction; # <--------------- 注释后为短事务
    while @count > 0
        do
            insert into test_idxquery_while_longtx_running(c, nonce_no_idx)
            VALUES (md5(UNIX_TIMESTAMP() + @count), md5(rand() + @count)),
                   (md5(UNIX_TIMESTAMP() + @count + 1), md5(rand() + @count + 1)),
                   (md5(UNIX_TIMESTAMP() + @count + 2), md5(rand() + @count + 2)),
                   (md5(UNIX_TIMESTAMP() + @count + 3), md5(rand() + @count + 3)),
                   (md5(UNIX_TIMESTAMP() + @count + 4), md5(rand() + @count + 4))
            on duplicate key update c=md5(concat(c, rand() * 1000));
set @count = @count - 1;
end while;
#     commit; # <---------------
end;

-- 50W条数据垫底
select count(1) from test_idxquery_while_longtx_running;

call test_idxquery_while_longtx_running_1();
