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

blog.controller('MainCtrl', function($scope, $timeout){
  var text1 = function() {
     $scope.text1= "DRIFTERS";
   }
   var text2 = function() {
     $scope.text2= "RACERS";
   }
   var text3 = function() {
     $scope.text3= "ADVENTURERS";
   }
   var text4 = function() {
     $scope.text4= "YOUR BLOG AWAITS...";
   }

  $timeout(text1, 500);
  $timeout(text2, 1000);
  $timeout(text3, 1500);
  $timeout(text4, 2000);
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
blog.controller('BlogController', ['$http', '$window', function($http, $window){

   var blog = this;
   blog.title = "Blogs";

   blog.posts = {};
  //  $http.get('https://s3-us-west-2.amazonaws.com/s.cdpn.io/110131/posts_1.json').success(function(data){
  //    blog.posts = data;
  //  });
  $http.get('/blogs').success(function(data) {
    blog.posts = data;
  });

   blog.tab = 'blog';

   blog.selectTab = function(setTab){
     blog.tab = setTab;
     console.log(blog.tab)
   };

   blog.isSelected = function(checkTab){
     return blog.tab === checkTab;
   };

   blog.post = {};

   blog.addPost = function(){
    var uniqueid = (Math.random() * 1000).toString();
    $http.post('/blogs', {UniqueId : uniqueid, Title: blog.post.title,
      Body: blog.post.body, Author: blog.post.author, Comments: [], Likes: 0,
      CreatedOn: Date.now()}).
    error(logError).
    success(function(){
      $window.location.href="/user";
    });
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

 blog.controller('UserBlogController', ['$http', '$window', function($http, $window){

    var blog = this;
    blog.title = "Blogs";

    blog.posts = {};
   //  $http.get('https://s3-us-west-2.amazonaws.com/s.cdpn.io/110131/posts_1.json').success(function(data){
   //    blog.posts = data;
   //  });
   $http.get('/user').success(function(data) {
     console.log(data);
     blog.posts = data;
   });

    blog.tab = 'blog';

    blog.selectTab = function(setTab){
      blog.tab = setTab;
      console.log(blog.tab)
    };

    blog.isSelected = function(checkTab){
      return blog.tab === checkTab;
    };

    blog.post = {};

    blog.addPost = function(){
     var uniqueid = (Math.random() * 1000).toString();
     $http.post('/user', {UniqueId : uniqueid, Title: blog.post.title,
       Body: blog.post.body, Author: blog.post.author, Comments: [], Likes: 0,
       CreatedOn: Date.now()}).
     error(logError).
     success(function(){
       $window.location.href="/blogs";
     });
     };

  }]);

  blog.controller('UserCommentController', function(){
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
