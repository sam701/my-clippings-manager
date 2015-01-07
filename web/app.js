angular.module('MyClippings', ['ui.router', 'angularFileUpload'])
.config(function($stateProvider, $urlRouterProvider){
  $urlRouterProvider.otherwise('/books');
  $stateProvider
    .state('bookList', {
      url: '/books',
      templateUrl: 'books.html',
      controller: 'BooksCtrl'
    })
    .state('clippings', {
      url: '/books/:bookId',
      templateUrl: 'clippings.html',
      controller: 'ClippingsCtrl'
    })
	.state('uploads', {
		url: '/uploads',
		templateUrl: 'uploads.html',
		controller: 'UploadsCtrl'
	})
	;
})
.controller('BooksCtrl', function($scope, $http) {
    $http.get("/books").success(function(data){
      $scope.books = data;
    });
})
.controller('ClippingsCtrl', function($scope, $http, $stateParams) {
    $http.get("/books/" + $stateParams.bookId).success(function(data){
      $scope.book = {
        clippings: data
      };
    });
})
.controller('UploadsCtrl', function($scope, FileUploader){
	$scope.uploader = new FileUploader({
		url: '/upload',
		autoUpload: true
	});
})
;
