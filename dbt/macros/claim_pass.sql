{% macro claim_pass(value_expr, threshold_expr, operator_expr) %}
case
  when {{ threshold_expr }} is null or cast({{ threshold_expr }} as varchar) = '' then null
  when {{ operator_expr }} = '>=' and {{ value_expr }} >= try_cast(cast({{ threshold_expr }} as varchar) as double) then true
  when {{ operator_expr }} = '<=' and {{ value_expr }} <= try_cast(cast({{ threshold_expr }} as varchar) as double) then true
  when {{ operator_expr }} = '==' and {{ value_expr }} = try_cast(cast({{ threshold_expr }} as varchar) as double) then true
  when {{ operator_expr }} = '>'  and {{ value_expr }} >  try_cast(cast({{ threshold_expr }} as varchar) as double) then true
  when {{ operator_expr }} = '<'  and {{ value_expr }} <  try_cast(cast({{ threshold_expr }} as varchar) as double) then true
  else false
end
{% endmacro %}
