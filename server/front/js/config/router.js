app.config(['$stateProvider','$urlRouterProvider', '$locationProvider', function($stateProvider, $urlRouterProvider, $locationProvider) {
	$urlRouterProvider

	.otherwise('/dashboard');

	$stateProvider

	.state('login', {	
		url: '/login',
		templateUrl : 'templates/login.html',
		controller: 'LoginCtrl'
	})

	.state('dashboard', {
		url: '/dashboard',
		templateUrl: 'templates/dashboard.html',
		controller: 'AppCtrl'
	})

	.state('dashboard.relations', {
		url: '/relations',
		templateUrl: 'templates/relations.html',
		controller: 'RelCtrl',
		controllerAs: 'rel'
	})

	.state('dashboard.transactions', {
		url: '/transactions/{id}',
		templateUrl: 'templates/transactions.html',
		controller: 'TxnCtrl',
		controllerAs: 'txn'
	})



	//$locationProvider.html5Mode(true);

}])