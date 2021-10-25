#include <reg52.h>
#include <intrins.h>

sbit pcf8591_sda = P3^3;
sbit pcf8591_scl = P3^2;

#define delayNOP(); {_nop_();_nop_();_nop_();_nop_();};

#define	PCF8591_WRITE	0x90
#define	PCF8591_READ 	0x91

uchar receivebuf[4];

bit bdata SystemError;

void iic_start(void)
{
	SDA = 1;     
	SCL = 1;
	delayNOP(); 
   	SDA = 0;
	delayNOP();
    SCL = 0;
}

void iic_stop(void)
{  	
	SDA = 0;
	SCL = 1;
	delayNOP();
	SDA = 1;
	delayNOP();
    SCL = 0;
}

 void iicInit(void)
   {
   SCL = 0;
   iic_stop();	
   }  

void slave_ACK(void)
{
	SDA = 0;  
	SCL = 1;
	delayNOP();	
	SCL = 0;
}

void slave_NOACK(void)
{ 
	SDA = 1;
	SCL = 1;
	delayNOP();
	SDA = 0;
    SCL = 0;  
}

void check_ACK(void)
{ 	    
    SDA = 1;
	SCL = 1;
	F0 = 0;
	delayNOP();   
	if(SDA == 1)
    F0 = 1;
   	SCL = 0;
}

void IICSendByte(uchar ch)
 
{
  	unsigned char idata n=8;

	while(n--)
	{ 
	if((ch&0x80) == 0x80)
	   {
	 		SDA = 1
			SCL = 1;
		    delayNOP();
			SCL = 0; 
	   }
		else
		{  
			SDA = 0;
			SCL = 1;
			delayNOP();
		  	SCL = 0;
		}
		ch = ch<<1;
	}
}

uchar IICreceiveByte(void)
{
	uchar idata n=8;
	uchar tdata=0;
	while(n--)
	{
	   SDA = 1;
	   SCL = 1;
	   tdata =tdata<<1;
	   	if(SDA == 1)
		  tdata = tdata|0x01;
		else 
		  tdata = tdata&0xfe;
	   SCL = 0;
	 }

	 return(tdata);
}

void DAC_PCF8591(uchar controlbyte,uchar w_data)
{    
	
	iic_start();
	delayNOP();

	IICSendByte(PCF8591_WRITE);
	check_ACK();
    if(F0 == 1)
	 { 
		SystemError = 1;
		return;
     }
    IICSendByte(controlbyte&0x77); 
	check_ACK();
    if(F0 == 1)
	 { 
		SystemError = 1;
		return;
	 }
    IICSendByte(w_data);
	check_ACK();
    if(F0 == 1)
	 { 
		SystemError = 1;
    	return;
	 }
	iic_stop();
	delayNOP();
	delayNOP();
	delayNOP();
	delayNOP();	
}

void ADC_PCF8591(uchar controlbyte)
{ 
    uchar idata receive_da,i=0;

	iic_start();

	IICSendByte(PCF8591_WRITE);
	check_ACK();
	if(F0 == 1)
	{
		SystemError = 1;
		return;
	}

	IICSendByte(controlbyte);
	check_ACK();
	if(F0 == 1)
	{
		SystemError = 1;
		return;
	}

    iic_start();
   	IICSendByte(PCF8591_READ);
	check_ACK();
	if(F0 == 1)
	{
		SystemError = 1;
		return;
	}
	 
    IICreceiveByte();
    slave_ACK();

	while(i<4)
	{  
	  receive_da=IICreceiveByte();
	  receivebuf[i++]=receive_da;
	  slave_ACK();
	}
	slave_NOACK();
	iic_stop();
}

void main(void)
{
    uchar i,l;
    delay(10);
    i = 0;
    while(dis4[i] != '\0')
     {                         
       lcd_wdat(dis4[i]);
       i++;
     }

    lcd_pos(0x40);
    i = 0;
    while(dis5[i] != '\0')
    {
        lcd_wdat(dis5[i]);
        i++;
    }

  while(1)
  {
	iicInit();
    ADC_PCF8591(0x04);

	if(SystemError == 1)
	  {
	  	iicInit();				  
	    ADC_PCF8591(0x04);
	   }   
	
	for(l=0;l<4;l++)	
	 {
	  show_value(receivebuf[0]);       
	    lcd_pos(0x02);        
        lcd_wdat(dis[2]);
        lcd_pos(0x04);
        lcd_wdat(dis[1]);
        lcd_pos(0x05);        
        lcd_wdat(dis[0]);

      show_value(receivebuf[1]);
	    lcd_pos(0x0b);
        lcd_wdat(dis[2]);
        lcd_pos(0x0d);       
        lcd_wdat(dis[1]);
        lcd_pos(0x0e);
        lcd_wdat(dis[0]);

	  show_value(receivebuf[2]);
	    lcd_pos(0x42);
        lcd_wdat(dis[2]);
        lcd_pos(0x44);          
        lcd_wdat(dis[1]); 
        lcd_pos(0x45);          
        lcd_wdat(dis[0]);

      show_value(receivebuf[3]);	 
	    lcd_pos(0x4b);   
        lcd_wdat(dis[2]);
        lcd_pos(0x4d);           
        lcd_wdat(dis[1]); 
        lcd_pos(0x4e);          
        lcd_wdat(dis[0]);

	  iicInit();
      DAC_PCF8591(0x40,receivebuf[0]);

	   	if(SystemError == 1)
	    {
	  	 iicInit();
		 DAC_PCF8591(0x40,receivebuf[0]);
	    }	        
	 }

   }
}
