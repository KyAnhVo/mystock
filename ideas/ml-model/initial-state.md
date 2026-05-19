# Initial State

## High level

Suppose we have a set of portfolios: 
```
P ⊆  {p | p ⊆ {(s, v) | s ∈ S and v ∈ R+}} # v is value in that stock
S := {s | s is a stock in US}
```

Consider the following functions:
- `v: S × R -> R^n`: converts a stock `s` at some time `t` to a vector embedding `v(s, t) ∈ R^n`
- `g: P × R -> R^n × R`: generates a stock vector and dollar value in that stock 
  from a given portfolio

We generate a portfolio as following:

```
proc generate_portfolio(t):
  let portfolio := {}
  let value_total := <some total portfolio value>
  while value_total > 0:
    let stock_vector, value := g(portfolio, v_total)
    let stock := choose_stock(stock_vector)
    portfolio := portfolio ∪ {(stock, value)}
    value_total := value_total - min(value, value_total)
  return portfolio
```

where choose_stock is defined as the following:

```
proc choose_stock(stock_vector, t):
  let stock_probability := []
  let distribution := normal distribution where mean is stock_vector and some constant covariance matrix
  for each remaining stock s:
    append distribution(v(s, t)) to stock_probability
  let n_p := normalized stock_probability
  choose a stock from n_p
```

We necessitates that:
- `∀s1,s2,s3 ∈ S: ∀t ∈ R+: (|v(s1, t) - v(s2, t)| < |v(s1, t) - v(s3, t)| <=> similarity(s1, s2) > similarity(s1, s3))`,
  for some similarity heuristic.
