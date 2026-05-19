# MyStock

A stock web application with personalized suggestions for stocks to buy (with AI).

## Idea

- Company stock viewer
- Stock suggestions based on personal portfolio

## How do we suggest stocks for you?

Consider some portfolio. A stock that fits a portfolio must have a combination of expanding
the portfolio's breadth (variety), strong stock, etc. We build a neural network system to 
evaluate a stock's S's fit to a portfolio P. Essentially, the function built here is
`f: S x pow(S) -> [0, 1]` where the input is a stock S and a portfolio P which is a subset
of all stocks, returning a grade from 0 to 1 of that S wrt P.

Our world will be modeled as a game, where it is player against market game. We consider the following
game with actions: `<buy|sell> <stock>` or `do nothing`, and we model our reward wrt some index fund.
