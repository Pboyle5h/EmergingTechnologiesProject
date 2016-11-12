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

// blog.controller('MainCtrl', function($scope){
//   $scope.isAuth = function(){
//     $scope.isAuth
//   }
// })
blog.controller('RegisterCtrl', function($scope, $http, $window){
  //console.log("called")
  $scope.register = function() {
    $http.post('/register', {Name: $scope.name, Username: $scope.username,
                     Email: $scope.email, Password: $scope.password}).
      error(logError).
      success(function(data) {
        $window.location.href="/";
      });
  };
});

blog.controller('LoginCtrl', function($scope, $http, $window){
  $scope.login = function(){
    $http.post('/login', {Username : $scope.username, Password : $scope.password}).
      error(function(){
        logError;
      }).
      success(function(){
        //$cookieStore.put('access_token', data.Access_token);
        //$cookieStore.put('refresh_token', data.Refresh_token);
        $window.location.href="/user";
      });
  };
});

var logError = function(data, status) {
   console.log('code '+status+': '+data);
 };
