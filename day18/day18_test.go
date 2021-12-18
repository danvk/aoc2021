package main

import "testing"

func TestExplode(t *testing.T) {
	p, _ := ParsePair("[3,2]")
	r, expL, expR, exploded := p.Explode(3)
	if r.String() != p.String() || expL != nil || expR != nil || exploded {
		t.Errorf("[3,2].Explode(3) == %s, %v, %v, %v want [3,2], nil, nil, false",
			r, expL, expR, exploded,
		)
	}

	r, expL, expR, exploded = p.Explode(4)
	if r.String() != "0" || *expL != 3 || *expR != 2 || !exploded {
		t.Errorf("[3,2].Explode(4) == %s, %v, %v, %v want 0, 3, 2, true",
			r, expL, expR, exploded,
		)
	}

	p, _ = ParsePair("[4,[3,2]]")
	r, expL, expR, exploded = p.Explode(3)
	if r.String() != "[7,0]" || expL != nil || expR == nil || *expR != 2 || !exploded {
		t.Errorf("%s.Explode(3) == %s, %v, %v, %v want [7, 0], nil, 2, true",
			p, r, expL, expR, exploded,
		)
	}

	p, _ = ParsePair("[5,[4,[3,2]]]")
	r, expL, expR, exploded = p.Explode(2)
	if r.String() != "[5,[7,0]]" || expL != nil || *expR != 2 || !exploded {
		t.Errorf("%s.Explode(2) == %s, %v, %v, %v want [5,[7,0]], nil, 2, true",
			p, r, expL, expR, exploded,
		)
	}
}

func TestSplit(t *testing.T) {
	p, _ := ParsePair("10")
	r, split := p.Split()
	if !split || r.String() != "[5,5]" {
		t.Errorf("%s.Split() = %s, %v want [5,5], true", p, r, split)
	}

	p, _ = ParsePair("11")
	r, split = p.Split()
	if !split || r.String() != "[5,6]" {
		t.Errorf("%s.Split() = %s, %v want [5,6], true", p, r, split)
	}

	p, _ = ParsePair("9")
	r, split = p.Split()
	if split || r.String() != "9" {
		t.Errorf("%s.Split() = %s, %v want 9, false", p, r, split)
	}
}
