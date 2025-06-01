$scoreboard players set $FUNCEXIT $(namespace) 0
$execute store result storage $(namespace):data vars.$CALL int 1 run scoreboard players get $CALL $(namespace)
$data modify storage $(namespace):data $(args).__call__ set from storage $(namespace):data vars.$CALL
$scoreboard players add $CALL $(namespace) 1
$execute store success score $FUNCEXIT $(namespace) run function $(function_namespace):$(function) with storage $(namespace):data $(args)
$execute unless score $FUNCEXIT $(namespace) matches 1 run function mcb:zzz/report {text:{text:'Function $(function_namespace):$(function) failed during execution',color:red,italic:true}}
return 0
