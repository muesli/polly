import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

const { inject: { service } } = Ember;

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    currentUser: service('current-user'),

    model() {
       return Ember.RSVP.hash({
           budget: this.store.findAll('budget')
       });
     },

     setupController(controller, models) {
       controller.set('budget', models.budget);
       controller.set('contact', this.get('currentUser').get('user').get('email'));
       controller.set('errorMessage', '');
       controller.set('responseMessage', '');
       controller.set('progressMessage', '');
     }
});
