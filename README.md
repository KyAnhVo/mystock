# MyStock

A stock web application with personalized suggestions for stocks to buy (with AI).

## Idea

- Company stock viewer
- Stock suggestions based on personal portfolio

## How do we suggest stocks for you?

### High level
We model the portfolio management job as a player-v-market game. We consider the following
game with actions: `<buy|sell> <stock>` or `do nothing` each trading day, and we compute
reward relative to various index fund benchmarks.

### Details: DQN
We train a DQN-based (e.g. DQN, DDQN, Dueling DQN, Noisy DQN, etc.) based on historical data,
simulating daily trading episodes. The game's initial state will be some form of randomization
of stock values (not totally random, created with a generator that is trained from example portfolios).

### Formal defition
Formally, we approximate a function`f: S \times pow(S) -> [0, 1]` where `S: {s | s is a stock in US}`
and for each porfolio P, we find the stock S: `S = argmax_s(f(s, P))`
