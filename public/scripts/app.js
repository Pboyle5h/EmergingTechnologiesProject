// adapated from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
var blog = angular.module('blogApp', ['ngRoute']);

blog.config(function($routeProvider, $locationProvider){
  $routeProvider
  .when('/', {templateUrl: '/partials/main.html'})
  .when('/blogs', {templateUrl: '/partials/blogs.html'})
  .when('/register', {templateUrl: '/partials/register.html'})
  .when('/login', {templateUrl: '/partials/login.html'});

  $locationProvider.html5Mode(true); // takes the # out of the url
});
