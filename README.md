<p align="center">
  <img src="https://github.com/basebandit/gocash/blob/master/change.png" alt="Icon PNG" height="64">
  <h3 align="center">gocash</h3>
  <p align="center">Convert Currency Rates directly from your Terminal!<p>

</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/basebandit/gocash)](https://goreportcard.com/report/github.com/basebandit/gocash)  [![GitHub license](https://img.shields.io/github/license/basebandit/gocash)](https://github.com/basebandit/gocash/blob/master/LICENSE)

<p align="center"><img src="https://github.com/basebandit/gocash/blob/master/cash.gif" alt="Icon PNG"></p>

## Install

```
go get -u github.com/basebandit/gocash/cmd/cash
```
Copy the following files to the `.gocash` directory inside your home directory
- [currencies.json](https://github.com/basebandit/gocash/blob/master/currencies.json) file
- [config.json](https://github.com/basebandit/gocash/blob/master/config.json) file 

## Usage
 
> You will need a Fixer.io API key. Get it [here](https://fixer.io/signup/free) for free.  



```bash
	Usage
		$ cash <amount> <from> <to>  
		$ cash <options>  
	Examples
		$ cash 10 usd eur pln aud kes tzs ugx
```

## Available Currencies

See [currencies.json](https://github.com/basebandit/gocash/blob/master/currencies.json) file.

## ToDo:
- Memoization 
- Tests 

## Thanks:
- [xxczaki](https://twitter.com/dokwadratu) for an awesome [tool](https://github.com/xxczaki/cash-cli), that inspired me to write a similar one in golang.
- [Money.js](http://openexchangerates.github.io/money.js/) for inspiring me to write my own golang version of the same library.
- [Fixer.io](http://fixer.io/) for providing awesome currency conversion API.
- <div>Icons made by <a href="https://www.flaticon.com/authors/smashicons" title="Smashicons">Smashicons</a> from <a href="https://www.flaticon.com/"                 title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/"                 title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

## License

[MIT](https://opensource.org/licenses/MIT)

