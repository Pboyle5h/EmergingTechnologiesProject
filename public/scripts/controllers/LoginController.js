//adapted from http://csfortheslothful.blogspot.ie/2014/11/try-web-development-in-go-with-slothful_30.html
blogApp.controller('LoginController', function($scope, $http){
  $scope.initialize = function(){
    $http.get('/login')
  };
});

blogApp.controller('RegistrationController', function($scope, $http){
  $scope.initialize = function(){
    $http.get('/register')
  };
});
