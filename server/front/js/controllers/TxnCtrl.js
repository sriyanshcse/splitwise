app.controller('TxnCtrl', ['$scope', '$stateParams', '$resource', '$rootScope' , function($scope, $stateParams, $resource, $rootScope){
	
	//console.log($stateParams.id);
	var vm = this;
	vm.transactions = [];
	$scope.id = $stateParams.id;
	
	if ($scope.id == -1) {
		var txn = $resource("/api/getAllTransactionHistory/:id");
		txn.query({id: $rootScope.getItem('id')}, 
			function (response) {
				vm.transactions = response;
				console.log(response);
			}
		);
	} else {
		var txn = $resource("/api/user/getTransactionHistory/:id1/:id2");
		txn.query({id1: $rootScope.getItem('id'), id2: $scope.id}, 
			function (response) {
				vm.transactions = response;
				console.log(response);
			}
		);
	}
	var month = new Array(); 
	month[0] = "Jan"; month[1] = "Feb";
	month[2] = "Mar";month[3] = "Apr";
	month[4] = "May";month[5] = "June";
	month[6] = "July";month[7] = "Aug";
	month[8] = "Sept";month[9] = "Oct";
	month[10] = "Nov";month[11] = "Dec";
	$scope.getDate = function (t) { 
		var date = t.split('T')[0].split('-');
		return date[2] + " " + month[parseInt(date[1])-1] + "," + date[0];
	}

	$scope.getAbs = function (amt) {
		return (amt > 0 ? amt : -amt);
	}
	$scope.getTime = function (t) {
		var time = t.split('T')[1].split('.')[0].split(':');
		var hr = parseInt(time[0]);
		hr = hr > 12 ? hr - 12 : hr;
		var k = hr >= 12 ? 'pm' : 'am';
		if (hr == 0) hr = 12;
		hr = hr.toString();
		if (hr.length == 1) {
			hr = '0' + hr;
		}
		return hr + ":" + time[1] + " " + k;
	}
}])