$data modify storage $(storage) $(res) set value $(a)
$execute store success storage $(storage) $(res) int 1 run data modify storage $(storage) $(res) set value $(b)
return 0