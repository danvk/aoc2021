# Day 24 working notes

Valid plates seem quite rare! I generated a million random ones and didn't find a single valid plate. I'm guessing this has to be calculating some mathematical property, but primality would produce too many valid plates.

mod: All mods are modulo 26.
div: All divs are divide by either 1 or 26.
inp: All inputs are to `w`, and only `inp` writes to `w`.

With all these mod 26s, I wonder if this is something alphabetical?

This sequence is repeated 14 times:

fn0:
		mul x 0   # x=0
		add x z   # x=z
		mod x 26  # x=(z%26)
		div z 1   # no-op

fn1:
    eql x w
    eql x 0   # x=1 if x=w else 0
    mul y 0
    add y 25  # y=25
    mul y x   # y=25*(x=1 if x=w else 0)
    add y 1   # y=1 + 25*(x=1 if x=w else 0)
    mul z y   # z = z * y
    mul y 0   # y = w
    add y w

		if x == w {
			x=0
		} else {
			z *= 26
			x=1
		}
		y = w

fn2:
		mul x 0
		add x z
		mod x 26  # x = z % 26
		div z 26  # z /= 26

The processing of each digit involves two numbers, and maybe a divide by 26:

 D0:      +14, +12
 D1:      +11,  +8
 D2:      +11,  +7
 D3:      +14,  +4               9
 D4: /26, -11,  +4  #  d5=d4-7   2  d5=(d4+4)%26-11=d4-7
 D5:      +12,  +1
 D6: /26,  -1, +10  #  d6=d5     9  d6=(d5+1)%26-1 = d5
 D7:      +10,  +8
 D8: /26,  -3, +12  # d8=d7+5       d8=(d7+8)%26-3 = d7+5
 D9: /26,  -4, +10  # d9=d8+8    54 d9=(d8+12)%26-4 = d8+8
D10: /26, -13, +15  # d10=d9-3   21 d10=(d9+10)%26-13 = d9-3
D11: /26,  -8,  +4  # d11=d10+7  98 d11=(d10+15)%26-8 = d10+7
D12:      +13, +10               9
D13: /26, -11,  +9  # d13=d12-1  8

The second column of numbers is all positive.
When there's a `/26`, the previous input's contribution to `z` is eliminated, except insofar as it could cause an `x==w` match.
The values in the first column are all positive and >=10, unless there's a `/26`, in which case they're negative.

The numbers that contribute to z are all the inputs, values in the second column and multiplying by 26. So how can you get a zero? The `x==w` cases at most let you drop terms.

I think I want to hit all the x==w cases?

So d5=d4-7

Leaving the group for an input, we always have:

x = 0 or 1
y = 0 or input + step.b
z = ???
w = input

z carries from one step to the next
w gets reset (to the next input)
x gets reset
y gets reset

  x=z%26
	if step.div { z/=26 }
	x += step.a
	if x == w {
		x = 0
		y = 0
	} else {
		z = 26 * z + w + step.b
	}

With this loop nice and fast, running over 1,000,000,000 inputs to get a sample of valid plates was a cinch. This let me figure out the patterns between digits:

D0D1 2 3 4 5 6 7 8 910111213
 2 8 6 9 2 2 2 1 6 9 3 6 7 6 evaluates to zero!

 d5 != d4 - 7
 d6 == d5
 d8 == d7 + 5
 d9 != d8 + 8
d10 != d9 - 3
d11 != d10 + 7
d13 == d12 - 1

28692221693676 (too low)

28692221693698 (max d13)
28692991693698 (max d5, d6)
28692994993698 (max d7, d8)  not correct

28692994993698
14 380 9893 257231 9893 257228 9893 257230 9893 380 14 0 19 0

28692221693676
14 380 9893 257231 9893 257221 9893 257227 9893 380 14 0 17 0

 0 2
 1 8
 2 6
 3 9 a= 14 b= 4
 4 2 a=-11 b= 4 match
 5 2 a= 12 b= 1
 6 2 a= -1 b=10 match
 7 1 a= 10 b= 8
 8 6 a= -3 b=12 match
 9 9 a= -4 b=10 match
10 3 a=-13 b=15 match
11 6 a= -8 b= 4 match
12 8 a= 13 b=10
13 7 a=-11 b= 9 match

So YES, you have to hit all the `x==w` cases.

d5=d6 seems to be a hard and fast constraint

                     1 1 1 1
       3 4 5 6 7 8 9 0 1 2 3
[2 8 6 9 2 2 2 1 6 9 3 6 7 6] -> 0
[4 9 3 9 2 6 6 4 9 6 4 8 9 8] -> 0
[4 6 1 8 1 6 6 4 9 4 1 8 9 8] -> 0
[5 7 5 8 1 4 4 3 8 8 2 9 5 4] -> 0
[1 8 6 9 2 2 2 4 9 9 3 5 2 1] -> 0

d4=d3-7
d6=d5
d8=d7+5
d13=d12-1

d9 has no obvious relation to d8
d10 has no obvious relation to d9
d11 has no obvious relation to d10

So:
d0-d2 have some relationship to d3 (1000 possibilities)
d4=d3-7; 2 possibilities: 81, 92
d5=d6; these may as well be 9
d8=d7+5; these may as well be 49
d13=d12-1; these may as well be 98

This constrains the problem enough that iterating over all possible plates is quite easy.

57581443882998
57592443882998
57592993882998
57592994982998

59692994994998
16181994941598 = too high
16181111641521
