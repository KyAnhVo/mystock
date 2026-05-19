# Initial State

## High level

Suppose we have a set of portfolios: 
```
P ⊆  {p | p ⊆ {(s, v) | s ∈ S and v ∈ R+}} # v is value in that stock
S := {s | s is a stock in US}
```

Consider the following functions:
- `v: S -> R^n`: converts a stock `s` to a vector `v(s) ∈ R^n`
- `g: P × R -> R^n × R`: generates a stock vector and dollar value in that stock 
  from a given portfolio

We generate a portfolio as following:

```
let P := {}
let v_total := <some total portfolio value>
while v_total > 0:
  s, v := g(P, v_total)
  P := P ∪ {(s, v)}
  v_total := v_total - min(v, v_total)
return P
```
