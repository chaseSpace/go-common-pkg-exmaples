drop procedure if exists test_idxquery_while_longtx_running_2;
create procedure test_idxquery_while_longtx_running_2()
begin
    set @query_max_id = 10000;
    while @query_max_id > 0
        do
            set @qid = floor(rand() * 10000);
select nonce_no_idx from test_idxquery_while_longtx_running where c = md5(@qid) limit 1 into @_;
set @query_max_id = @query_max_id - 1;
end while;
end;

call test_idxquery_while_longtx_running_2();

show processlist;
