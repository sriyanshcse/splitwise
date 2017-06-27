app.controller('AppCtrl', ['$scope', '$mdBottomSheet','$mdSidenav', '$mdDialog', '$rootScope', function($scope, $mdBottomSheet, $mdSidenav, $mdDialog, $rootScope){
  $scope.toggleSidenav = function(menuId) {
    $mdSidenav(menuId).toggle();
  };
  console.log("hello");
  

 	$scope.menu = [
    {
      link : '',
      title: 'Dashboard',
      icon: 'dashboard',
      state: 'dashboard.relations'
    },
    {
      link : '',
      title: 'Friends',
      icon: 'group',
      state: 'dashboard.relations'
    },
    {
      link : '',
      title: 'Notifications',
      icon: 'message',
      state: 'dashboard.notifications'
    }, 
    {
      link : '',
      title: 'Transactions',
      icon: 'description',
      state: 'dashboard.transactions({id:-1})'
    }
  ];

  $scope.admin = [
    {
      link : '',
      title: 'Trash',
      icon: 'delete'
    },
    {
      link : 'showListBottomSheet($event)',
      title: 'Settings',
      icon: 'settings'
    }
  ];
  $scope.search;
  $scope.$watch('search', function (newVal, oldVal) {
    $rootScope.search = newVal;
    console.log(newVal);
  });
  $scope.alert = '';
  $scope.showListBottomSheet = function($event) {
    $scope.alert = '';
    $mdBottomSheet.show({
      templateUrl: 'templates/bottomSheet.html',
      controller: 'ListBottomSheetCtrl',
      targetEvent: $event
    }).then(function(clickedItem) {
      $scope.alert = clickedItem.name + ' clicked!';
    });
  };
  
  $scope.showAdd = function(ev) {
    $mdDialog.show({
      controller: 'DialogController',
      templateUrl: 'templates/dialog.html',
      targetEvent: ev
    })
    .then(function() {
      $scope.alert = 'You said the information was';
    }, function() {
      $scope.alert = 'You cancelled the dialog.'	;
    });
  };
 $scope.showNameDialog = function() {
  $mdDialog.show({
      controller: 'nameDialogController',
      templateUrl: 'templates/nameDialog.html',
      clickOutsideToClose: false,
      escapeToClose: false
    })
    .then(function() {
      $scope.alert = 'You said the information was';
    }, function() {
      $scope.alert = 'You cancelled the dialog.'	;
    });
  }
  if (!$rootScope.getItem('firstName')) {
    console.log("showing Dialog");
    $scope.showNameDialog();
  }
}]);