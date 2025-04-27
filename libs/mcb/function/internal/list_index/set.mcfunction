$data modify storage $(storage) $(list)$(index) set from storage $(storage) $(value_path)
$tellraw @a {text:"$(storage)$(list)$(index) = $(value_path)"}