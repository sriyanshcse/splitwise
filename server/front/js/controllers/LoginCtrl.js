
angular.module ('Splitwise')
       .controller ('LoginCtrl', ['$scope', '$state', 'AuthenticationService', '$interval', '$resource', '$rootScope',LoginCtrl]);


function LoginCtrl ($scope, $state, AuthenticationService, $interval, $resource, $rootScope) {
    
    $scope.userOTP = "";
    $scope.userid;
    $scope.OTP;
    $scope.wrongOTP = false;
    var stop, sec = 120;
    $scope.mint = '02', $scope.secs = '00';
    $scope.isPaneShown = false;

    var timer = function () {
      stop = $interval (function () {
        sec--;
        if (!sec) $scope.reset ();
        var t = parseInt(sec / 60);
        $scope.mint = '0' + t.toString ();
        t = sec - t * 60;
        $scope.secs = '';
        if (t < 10) $scope.secs = '0';
        $scope.secs += t.toString (); 
        console.log ($scope.mint + " " + $scope.secs);
     }, 1000);
    }
    console.log($scope.userid);
    $scope.sendOTP = function () {
        if(!$scope.myForm.$valid) {
            return;
        }
        $scope.showPane();
       var genotp = $resource("/api/generateOTP");
       $scope.loading = true;
       genotp.save({Email: $scope.userid}, function(response) {
            $scope.loading = false;
            $scope.OTP = response.OTP;
            console.log($scope.OTP);
            timer ();
       });
       console.log("Hello");
    }

    // call api
    $scope.verifyOTP  = function () {
        console.log ($scope.userOTP);
        if ($scope.OTP == $scope.userOTP) {
            $scope.wrongOTP = false;
            console.log($scope.userid);
            AuthenticationService.login ($scope.userid);
            $interval.cancel (stop);
            stop = undefined;
            var checkuser = $resource("/api/checkUserExists/:eid");
            checkuser.get({eid : $scope.userid}, {}, function(response) {
            if(response.Exist) {
                $rootScope.setItem("id", response.Id);
                console.log("sadasda" + response.Name);
                var name = response.Name.split(' ');
                console.log(name);
                $rootScope.setItem('firstName', name[0]);
                $rootScope.setItem('lastName', name[1]);
            }
            $state.go('dashboard');
       });
            
        } else {
            $scope.wrongOTP = true;
            $scope.userOTP = '';
            //console.log ("Wrong OTP " + $scope.userOTP);
        }
    }
    
    $scope.reset = function () {
        if (angular.isDefined(stop)) {
            $interval.cancel (stop);
            stop = undefined;
            sec = 120, $scope.mint = '02', $scope.secs = '00';
            console.log("stopping $interval -stop")
        }
        $scope.OTP = '';
        $scope.wrongOTP = false;
        console.log ($scope.OTP);
    }
    $scope.showPane = function() {
        $scope.isPaneShown = true;
    };
    $scope.hidePane = function() {
        $scope.isPaneShown = false;
    };

    $scope.emailFormat = /^[a-z]+[a-z0-9._]+@[a-z]+\.[a-z.]{2,5}$/;
    
}