import Ember from 'ember';
import moment from 'moment';

export default Ember.Controller.extend({
    maxmicrobudget: 0,

    maxBudget: Ember.computed('maxmicrobudget', function() {
        this.store.query('budget', {
            month: moment().add(14, 'd').toDate().getMonth() + 1
        }).then((budget) => {
            this.set('maxmicrobudget', budget.objectAt(0).get('value'));
        });

        return this.get('maxmicrobudget');
    })
});
