$scoreboard players set $FUNCEXIT $(namespace) 0
$execute store success score $FUNCEXIT $(namespace) run function $(function_namespace):$(function) with storage $(namespace):data $(args)
$execute unless score $FUNCEXIT $(namespace) matches 1 run function mcb:zzz/report {text:{text:'Function $(function_namespace):$(function) failed during execution',color:red,italic:true}}
return 0
