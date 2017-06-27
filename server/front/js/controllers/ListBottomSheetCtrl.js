app.controller('ListBottomSheetCtrl', function($scope, $mdBottomSheet, AuthenticationService) {
  $scope.items = [
    { name: 'Share', icon: 'share' },
    { name: 'Upload', icon: 'upload' },
    { name: 'Copy', icon: 'copy' },
    { name: 'Logout', icon: 'print' },
  ];
  
  $scope.listItemClick = function($index) {
    var clickedItem = $scope.items[$index];
    $mdBottomSheet.hide(clickedItem);
    if (clickedItem.name === "Logout") {
      AuthenticationService.logout();
    }
  };
});