create or replace function public.get_paging_product(q character varying, page integer DEFAULT 1, pagesize integer DEFAULT 20)
 returns table(product_id bigint)
 language plpgsql
as $$
declare
	pattern varchar;
	_limit int;
	_offset int;
begin	
	pattern := concat('%', q, '%');
	_limit := pageSize;
	_offset := (page - 1) * _limit;

	return query select p.id::bigint
	from product p
	where p."name" like pattern
	offset _offset
	fetch next _limit rows only;
end; $$;

create or replace function public.get_total_count_product(q character varying)
 returns bigint
 language plpgsql
as $$
declare 
	pattern varchar;
	_result bigint;
begin 
	pattern := concat('%', q, '%');
	
	select count(*)::bigint as total_count
		into _result
		from product p
		where p."name" like pattern;

	return _result;	
end; $$;

-- DEMO
select * from get_paging_product('B', 1, 5);
select * from get_total_count_product('B') 
