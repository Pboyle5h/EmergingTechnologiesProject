// adapated from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
var blog = angular.module('blogApp', ['ngRoute']);

blog.config(function($routeProvider, $locationProvider){
  $routeProvider
  .when('/', {templateUrl: '/partials/main.html'})
  .when('/blogs', {templateUrl: '/partials/blogs.html'})
  .when('/register', {templateUrl: '/partials/register.html'})
  .when('/login', {templateUrl: '/partials/login.html'})
  .when('/user', {templateUrl: '/partials/user.html'});

  $locationProvider.html5Mode(true); // takes the # out of the url
});


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
        //console.log("error");
        $scope.invalidLogin = !$scope.invalidLogin;
      }).
      success(function(){
        //console.log("success");
        $window.location.href="/user";
      });
  };
});

// adapted from https://codepen.io/nickmoreton/pen/mgtLK
blog.controller('BlogController', ['$http', function($http){

   var blog = this;
   blog.title = "Blogs";

   blog.posts = {};
  //  $http.get('https://s3-us-west-2.amazonaws.com/s.cdpn.io/110131/posts_1.json').success(function(data){
  //    blog.posts = data;
  //  });
  $http.get('/blogs').success(function(data) {
    console.log(data);
    blog.posts = data;
  })

   blog.tab = 'blog';

   blog.selectTab = function(setTab){
     blog.tab = setTab;
     console.log(blog.tab);
   };

   blog.isSelected = function(checkTab){
     return blog.tab === checkTab;
   };

   blog.post = {};
   blog.addPost = function(){
     blog.post.createdOn = Date.now();
     blog.post.comments = [];
     blog.post.likes = 0;
     blog.posts.unshift(this.post);
     blog.tab = 0;
     blog.post ={};
   };

 }]);

 blog.controller('CommentController', function(){
   this.comment = {};
   this.addComment = function(post){
     this.comment.createdOn = Date.now();
     post.comments.push(this.comment);
     this.comment ={};
   };
 });

var logError = function(data, status) {
   console.log('code '+status+': '+data);
 };
