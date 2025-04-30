$scoreboard players set $FUNCEXIT $(namespace) 0
$data modify storage $(namespace):data args.$(function).__call__ set from storage $(namespace):data vars.$CALL
$scoreboard players add $CALL $(namespace) 1
$function $(namespace):$(function) with storage $(args)
$execute unless score $FUNCEXIT $(namespace) matches 1 run function mcb:error {text:'Function $(function) failed during execution'}
$data modify storage $(namespace):data vars.$(ret) set from storage $(namespace):data vars.$RET
return 1
