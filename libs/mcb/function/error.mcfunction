$tellraw @a[tag=mcblog] {text:'$(text)',italic:true,color:red}
execute as @a[tag=mcblog] at @s run playsound minecraft:block.note_block.bass master @s ~ ~ ~ 1 0
return 1
