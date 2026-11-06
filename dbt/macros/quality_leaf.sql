{% macro quality_leaf(applicable_expr, value_expr, comparison_key) %}
case
  when not ({{ applicable_expr }}) then json_object('status', 'not_applicable')
  when {{ value_expr }} is null then json_object('status', 'unknown')
  else json_object(
    'status', 'measured',
    'value', {{ value_expr }},
    'attribution', 'model_only',
    'comparison_key', '{{ comparison_key }}'
  )
end
{% endmacro %}
