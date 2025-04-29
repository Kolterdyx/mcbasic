$scoreboard players set $callstatus $(namespace) 0
$execute store success score $callstatus $(namespace) run function $(function) with storage $(args)
$execute unless score $callstatus $(namespace) matches 1 run function mcb:error {text:'Function $(function) failed during execution'}
return 1
