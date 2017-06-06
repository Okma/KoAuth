/**
 * Created by Carl on 6/6/2017.
 */

// Inject jQuery.
let script = document.createElement('script');
script.src = 'http://code.jquery.com/jquery-3.2.1.min.js';
script.type = 'text/javascript';
script.integrity = "sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4=";
script.crossOrigin = "anonymous";
document.getElementsByTagName('head')[0].appendChild(script);

// Inject sign up js.
script = document.createElement('script');
script.src = 'js/signup.js';
script.type = 'text/javascript';
document.getElementsByTagName('head')[0].appendChild(script);

// Inject usage js.
/*
js = document.createElement('js');
js.src = '/js/signup.js';
js.type = 'text/javascript';
document.getElementsByTagName('head')[0].appendChild(js);*/
