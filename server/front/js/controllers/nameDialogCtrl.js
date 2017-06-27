
app.controller('nameDialogController', ['$scope', '$rootScope', '$mdDialog', '$state', '$resource', function ($scope, $rootScope, $mdDialog, $state, $resource) {
  $scope.hide = function() {
    $mdDialog.hide();
  };
  $scope.firstName = ""; $scope.lastName = "";
  $scope.cancel = function() {
    $mdDialog.cancel();
  };
  $scope.answer = function() {
  	if ($scope.firstName === "" || $scope.lastName === "") return;
    $rootScope.setItem("firstName", $scope.firstName);
    $rootScope.setItem("lastName", $scope.lastName);
    var addUser = $resource("/api/addUser");
    console.log("Adding User");
    addUser.save({Email : $rootScope.getItem("username"), Name: $scope.firstName + " " + $scope.lastName}, function (response) {
      console.log(response);
      $rootScope.setItem("id", response.Id);
    });
    $mdDialog.hide();
    $state.go('dashboard.relations');
  };
}]);