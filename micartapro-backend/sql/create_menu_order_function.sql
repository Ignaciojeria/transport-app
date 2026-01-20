create or replace function public.create_menu_order(
  p_menu_id uuid,
  p_event_payload jsonb,
  p_event_type text
)
returns table(order_number int)
language plpgsql
as $$
declare
  v_order_number int;
begin
  -- Genera correlativo por menú de forma atómica
  insert into public.menu_order_counters (menu_id, next_value)
  values (p_menu_id, 2)
  on conflict (menu_id)
  do update set next_value = public.menu_order_counters.next_value + 1
  returning public.menu_order_counters.next_value - 1
  into v_order_number;

  insert into public.menu_orders (
    menu_id,
    order_number,
    event_payload,
    event_type
  ) values (
    p_menu_id,
    v_order_number,
    p_event_payload,
    p_event_type
  );

  return query select v_order_number;
end;
$$;