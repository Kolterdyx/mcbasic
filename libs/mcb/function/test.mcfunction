$execute store success score test example run function $(function)
$execute if score test example matches 1 run tellraw @a[tag=mcblog] {text:'Function $(function) executed successfully',color:green}
$execute unless score test example matches 1 run tellraw @a[tag=mcblog] {text:'Function $(function) failed during execution',color:red}