import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendLaqpayComponent } from './send-laqpay.component';

describe('SendLaqpayComponent', () => {
  let component: SendLaqpayComponent;
  let fixture: ComponentFixture<SendLaqpayComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendLaqpayComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendLaqpayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
