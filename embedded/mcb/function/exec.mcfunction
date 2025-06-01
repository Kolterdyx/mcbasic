$scoreboard players set $FUNCEXIT $(namespace) 0
$execute store success score $FUNCEXIT $(namespace) run $(command)
$execute unless score $FUNCEXIT $(namespace) matches 1 run function mcb:zzz/report {text:{text:'Command $(command) failed during execution',color:red,italic:true}}
return 0
