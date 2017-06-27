
app.controller('BillDialogController', ['$scope', '$mdDialog', '$resource', '$rootScope', 'id', function ($scope, $mdDialog, $resource, $rootScope, id) {
  $scope.tname="";$scope.value="";
  $scope.hide = function() {
    $mdDialog.hide();
  };
  $scope.cancel = function() {
    $mdDialog.cancel();
  };
  $scope.id = id;
  console.log("bill dialog: " + $scope.id);
  $scope.answer = function() {
    if ($scope.tname === "" || $scope.value === "" || parseInt($scope.value) <= 0) return;
  	var addBill = $resource("/api/user/addBill/:id");
  	console.log("Adding Bill for " + $rootScope.getItem("id") + " " + $scope.id);
  	addBill.save({id: $rootScope.getItem("id")}, {Name: $scope.tname, To: $scope.id, Value: parseInt($scope.value)}, function (response) {
  		 console.log($rootScope.friends);
      //console.log(response);
  	});	
    $mdDialog.hide();
  };
}]);