import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    activate: function() {
        this._super();
        Ember.$('grantsFinished').button('toggle');
    },

    model() {
       return Ember.RSVP.hash({
         finishedProposals: this.store.query('proposal', {
             ended: true
         })
       });
     },

     setupController(controller, models) {
       controller.set('finishedProposals', models.finishedProposals);
     }
});
