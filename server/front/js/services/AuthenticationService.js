angular.module('Splitwise')
    .service('AuthenticationService', ['$window', '$state', function($window, $state) {
  this.checkIfLoggedIn = function(callback) {
  
        if ($window.localStorage['username']) {
            callback (true);
        } else {
            callback (false);
        }         
  }
  
  
  this.login = function(username) {
            // store username and token in local storage to keep user logged in between page refreshes
            $window.localStorage['username'] = username; 
  }
  
  this.logout = function () {
      console.log ("logging out");
      $window.localStorage.clear();
      $state.go($state.current, {}, {reload: true});
  }
                    
}]);