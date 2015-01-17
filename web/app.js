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
		var out = _.map(data, function(val, key){
			val.id = key;
			val.ClippingsTimes.LastFormatted = moment(val.ClippingsTimes.Last*1000).format('YYYY-MM-DD');
			return val;
		});
		out = _.sortBy(out, function(n){return -n.ClippingsTimes.Last;});
      $scope.books = out;
    });
})
.controller('ClippingsCtrl', function($scope, $http, $stateParams) {
    $http.get("/books/" + $stateParams.bookId).success(function(data){
      $scope.book = data;
    });
})
.controller('UploadsCtrl', function($scope, $http, FileUploader){
	$scope.uploader = new FileUploader({
		url: '/upload',
		autoUpload: true
	});
	$scope.uploader.onSuccessItem = function(item, res) {
		var it = _.find($scope.uploads, {Id: res.Id});
		if(it){
			it.cl = "info";
		}else{
			res.cl = "info";
			$scope.uploads.push(res);
		}
	};
	
	$scope.uploads = [];
	$http.get('/uploadIndex').success(function(data){
		$scope.uploads = data || [];
	});
})
.filter('bytes', function() {
	return function(bytes, precision) {
		if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
		if (typeof precision === 'undefined') precision = 1;
		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
			number = Math.floor(Math.log(bytes) / Math.log(1024));
		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
	}
})
;
