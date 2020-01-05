var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
const cacheControl = require('express-cache-controller');
require('dotenv').config();

var indexRouter = require('./routes/index');

var app = express();
app.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "https://radio.sun-yryr.com");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

app.use(cacheControl({ maxAge: 0 }));
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use('/', indexRouter);

module.exports = app;
