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
proc generate_portfolio():
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
proc choose_stock(stock_vector):
  let stock_probability := []
  let distribution := normal distribution where mean is stock_vector and some constant covariance matrix
  for each remaining stock s:
    append distribution(v(s)) to stock_probability
  let n_p := normalized stock_probability
  choose a stock from n_p
```

**Note**: the covariance matrix needs to be chosen very carefully.
