app.controller('RelCtrl', ['$scope', '$rootScope', '$resource', '$mdDialog', '$state', function($scope, $rootScope, $resource, $mdDialog, $state) {
	$rootScope.friends = [];
  $rootScope.Mid = new Object();
  $scope.getState = function(id) {
    return 'dashboard.transactions({id:' + id + '})';
  }
  var vm = this;
  $scope.search = $rootScope.search;
  vm.myFriends = $rootScope.friends;

	/*$scope.hotels = [
      {
        id : 1,
        name : 'TGB', 
        city : 'Surat', 
        Value : '4500',
        state : 'dashboard.transactions({id: 1})'
      },
      {
        id : 2,
        name : 'Raddison-Blu',
        city : 'Surat',
        Value : '5500',
        state : 'dashboard.transactions({id: 2})'
      },
      {
        id : 2,
        name : 'Gateway Hotel',
        city : 'Surat',
        price : '5500',
        state : 'dashboard.transactions({id: 3})'
      },
      {
        id : 2,
        name : 'Sheraton',
        city : 'Surat',
        Value : '7500',
        state : 'dashboard.transactions({id: 4})'
      },
      {
        id : 2,
        name : 'Oberoi Hotel',
        city : 'Surat',
        Value : '7000',
        state : 'dashboard.transactions({id: 5})'
      },
    ];*/


    var getFriends = $resource("/api/getFriends/:id");
    var response = getFriends.query({id : $rootScope.getItem('id')}); 
    response.$promise.then(function (response) {
        console.log(response);
        $rootScope.friends = response;
        vm.myFriends = response;
        $rootScope.$apply();      
        console.log($rootScope.friends);
    });   

    $scope.getAbs = function (amt) {
      return amt > 0 ? amt : -amt;
    }


    $scope.settlePayment = function (index) {
      if(!$rootScope.friends[index].Value) return;
      var id = $rootScope.friends[index].Id;
      var settle = $resource("/api/settlePayment/:id1/:id2", 
        { 
          id1: $rootScope.getItem('id'),
          id2: id
        }, {
          update : {
            method: 'PUT'
          }
        });
      settle.update(function() {
        console.log("updated");
      }); 
    }

    $scope.showBillDialog = function(ev, index) {
    console.log("this is index " + $rootScope.friends[index].Id);
    $mdDialog.show({
      controller: 'BillDialogController',
      templateUrl: 'templates/billDialog.html',
      targetEvent: ev,
      locals: {
        id: $rootScope.friends[index].Id
      },
      controllerAs: 'ctrl'
    })
    .then(function() {
      $scope.alert = 'You said the information was';
    }, function() {
      $scope.alert = 'You cancelled the dialog.'  ;
    });
  };


}]);
