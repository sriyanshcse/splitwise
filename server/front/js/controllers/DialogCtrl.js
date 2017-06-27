
app.controller('DialogController', ['$scope', '$mdDialog', '$resource', '$rootScope', function ($scope, $mdDialog, $resource, $rootScope) {
  $scope.email;
  $scope.hide = function() {
    $mdDialog.hide();
  };
  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.answer = function() {
  	var addFriend = $resource("/api/addFriend/:id");
  	console.log("Adding friend " + $rootScope.getItem("id"));
  	addFriend.save({id: $rootScope.getItem("id")}, {Email : $scope.email}, function (response) {
  		console.log(response);
  	});
    $mdDialog.hide();
  };
}]);