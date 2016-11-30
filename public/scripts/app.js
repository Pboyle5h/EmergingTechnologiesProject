// adapated from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
var blog = angular.module('blogApp', ['ngRoute']);


blog.config(function($routeProvider, $locationProvider){
  $routeProvider
  .when('/', {templateUrl: '/partials/main.html'})
  .when('/blogs', {templateUrl: '/partials/blogs.html'})
  .when('/register', {templateUrl: '/partials/register.html'})
  .when('/login', {templateUrl: '/partials/login.html'})
  .when('/user', {templateUrl: '/partials/user.html'})
  .when('/logout', {templateUrl: '/partials/logout.html'})
  .when('/about', {templateUrl: '/partials/about.html'});

  $locationProvider.html5Mode(true); // takes the # out of the url
});

blog.run(function($rootScope){
  $rootScope.isAuth = false;
  $rootScope.username = "";
});

blog.controller('RegisterCtrl',function($scope, $http, $window){
  //console.log("called")
  $scope.register = function() {
    $http.post('/register', {Name: $scope.name, Username: $scope.username,
                     Email: $scope.email, Password: $scope.password}).
      error(logError).
      success(function(data) {
        $window.location.href="/login";
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

blog.controller('LoginCtrl', function($scope, $http, $location,authService){
  $scope.login = function(){
    authService.Login($scope.username, $scope.password, function(response, status){
        if(status == 200){
           //authService.setCredentials(response.username, response.password);
          //$window.location.href='/';
          $location.path("/");
        } else {
          console.log(response.status);
          $scope.invalidLogin = true;
        }
    });
  };
});

blog.controller('LogoutCtrl', function($scope, $http, $location, authService){
  $scope.logout = function(){
    authService.Logout();
  };

  $scope.changeMind = function(){
    $location.path("/");
  }
});

// adapted from https://codepen.io/nickmoreton/pen/mgtLK
blog.controller('BlogController', ['$http', '$window', function($http, $window){

   var blog = this;
   blog.title = "Blogs";

   $(window).keydown(function(event){
    if((event.which== 13) && ($(event.target)[0]!=$("textarea")[0])) {
      event.preventDefault();
      return false;
    }
  });

   blog.posts = {};
  //  $http.get('https://s3-us-west-2.amazonaws.com/s.cdpn.io/110131/posts_1.json').success(function(data){
  //    blog.posts = data;
  //  });
  $http.get('/blogs').success(function(response) {
    console.log(response);
    blog.posts = response;
  });

   blog.tab = 'blog';

   blog.selectTab = function(setTab){
     blog.tab = setTab;
     //console.log(blog.tab)
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
      $window.location.href="/";
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
      //console.log(blog.tab)
    };

    blog.isSelected = function(checkTab){
      return blog.tab === checkTab;
    };


    blog.post = {};
    blog.addPost = function(){
     var uniqueid = (Math.random() * 1000).toString();
     $http.post('/user', {UniqueId : uniqueid, Title: blog.post.title,
       Body: blog.post.body, Author: blog.post.author, Likes: 0,
       CreatedOn: Date.now(), Comments: []}).
     error(logError).
     success(function(){
       $window.location.href="/";
     });
     };

     blog.deletePost = function(post) {
    $http({
        url: '/user',
        method: 'DELETE',
        data: {
            UniqueId: post.uniqueid,
            Title: post.title,
            Body: post.body,
            Author: post.author,
            Likes: post.likes,
            CreatedOn: post.createon,
            Comments: post.comments
        },
        headers: {
            "Content-Type": "application/json;charset=utf-8"
        }
    }).success(function(res) {
        $window.location.href="/";
    }).error(logError);
    };

     blog.editPost = function(post){
       blog.editPost = true;
       blog.post = post;
     }
     blog.updatePost = function(){
       $scope.editing = true;
       $http.put('/user', {"body" : blog.post.body}).
       error(logError).
       success(function(data){
         blog.posts = data;
       })
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

blog.factory('authService', function($http, $rootScope, $location) {

  var service = {};
  //var username = "";
  service.Login = function(username, password, callback){
    $http.post('/login', {Username: username, Password: password}).
    success(function(response, status){
      console.log(response + " "  + status);
      console.log("username from service.login : " + username );
      //service.setCredentials(username, password);
      //service.setCredentials(username, password);
      $rootScope.isAuth = true;
      $rootScope.username = username;
      callback(response, status);
    });
  };

  service.Logout = function(){
    $http.post('/logout').
    success(function(response, status){
      $rootScope.isAuth = false;
      $rootScope.username = "";
      $location.path("/");
      //callback(response, status);
    });
  };

  return service;

});
