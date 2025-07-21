package packformat

type PackFormat int

// https://minecraft.wiki/w/Pack_format

const (
	V1_13_V1_14_4   PackFormat = 4
	V1_15_V1_16_1   PackFormat = 5
	V1_16_2_V1_16_5 PackFormat = 6
	V1_17_V1_17_1   PackFormat = 7
	V1_18_V1_18_1   PackFormat = 8
	V1_18_2         PackFormat = 9
	V1_19_V1_19_3   PackFormat = 10
	V1_19_4         PackFormat = 12
	V1_20_V1_20_1   PackFormat = 15
	V1_20_2         PackFormat = 18
	V1_20_3_V1_20_4 PackFormat = 26
	V1_20_5_V1_20_6 PackFormat = 41
	V1_21_V1_21_1   PackFormat = 48
	V1_21_2_V1_21_3 PackFormat = 57
	V1_21_4         PackFormat = 61
	V1_21_5         PackFormat = 71
	V1_21_6         PackFormat = 80
)
