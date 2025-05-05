$data modify storage example:vars $RET set value "$(from)"
execute store result storage example:vars $RET int 1 run data get storage example:vars $RET
return 0
