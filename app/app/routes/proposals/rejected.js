import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    activate: function() {
        this._super();
        Ember.$('grantsRejected').button('toggle');
    },

    model() {
       return Ember.RSVP.hash({
         rejectedProposals: this.store.query('proposal', {
             accepted: false,
             ended: true
         })
       });
     },

     setupController(controller, models) {
       controller.set('rejectedProposals', models.rejectedProposals);
     }
});
