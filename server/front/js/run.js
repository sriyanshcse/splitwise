angular.module ('Splitwise')
    .run (['$rootScope', '$state', '$stateParams', 'AuthenticationService', '$window',
         function ($rootScope, $state, $stateParams, AuthenticationService, $window) {
            var callCount = 0;
            $rootScope.getItem = function(item) {
                return $window.localStorage.getItem(item);
            }  
            $rootScope.search;
            $rootScope.setItem = function(item, value) {
                $window.localStorage[item] = value;
            }
            $rootScope.$on ('$stateChangeStart', function (event, to, toParams, from, fromParams) {
                AuthenticationService.checkIfLoggedIn (function (response) {
                    console.log (from.name + " " + to.name + " " + callCount++);
                    console.log (response);
                    if (!response && to.name != 'login') {
                        event.preventDefault();
                        $state.go ('login');
                    } else if (response && to.name == 'login') {
                       event.preventDefault ();
                    }
             });
        });
    }
]);