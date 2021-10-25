#include <reg52.h>

sbit 24c08_scl = P3^4;
sbit 24c08_sda = P3^5;

#define	OP_READ	0xa1
#define	OP_WRITE 0xa0

// delay function

// delay n(ms) function

void start()
// ��ʼλ
{
	SDA = 1;
   SCL = 1; 
	_nop_();
	_nop_();
	_nop_();
	_nop_();
	SDA = 0;
	_nop_();
	_nop_();
	_nop_();
	_nop_();
	SCL = 0;
}

void stop()

{
	SDA = 0;
	SCL = 1;
	_nop_();
	_nop_();
	_nop_();
	_nop_();
	SDA = 1;
	_nop_();
	_nop_();
	_nop_();
	_nop_();
	SDA=0;
	SCL=0;
}

unsigned char ReadData()
{
	unsigned char i;
	unsigned char x;
	for(i = 0; i < 8; i++)
	{
		SCL = 1;
		x<<=1;
		x|=(unsigned char)SDA;
		SCL = 0;
	}
	return(x);
}

bit WriteCurrent(unsigned char y)
{
	unsigned char i;
	bit ack_bit;
	for(i = 0; i < 8; i++)
	{
    	SDA = (bit)(y&0x80);
		_nop_();   
	   SCL = 1;      
   	_nop_(); 
	  _nop_();
		
	  	SCL = 0;
		y <<= 1;      
	}
	SDA = 1;
	_nop_(); 
	_nop_();
	SCL = 1;
	_nop_();
	_nop_();_nop_(); 
	_nop_(); 
	ack_bit = SDA;
	SCL = 0;
	return  ack_bit;
}

void WriteSet(unsigned char add, unsigned char dat)
{
	start();
	WriteCurrent(OP_WRITE);
	WriteCurrent(add);
	WriteCurrent(dat);
	stop();                
	delaynms(4);	       
}

unsigned char ReadCurrent()
{
	unsigned char x;
	start();
	WriteCurrent(OP_READ);  
	x=ReadData();  
	stop(); 
	return x;              
}

unsigned char ReadSet(unsigned char set_addr)
{
	start(); 
	WriteCurrent(OP_WRITE);      
	WriteCurrent(set_addr);       
	return(ReadCurrent());        
}


main(void)
{
   unsigned char i; 
   SDA = 1;         
	SCL = 1;  	       
   for(i = 0 ; i < 16; i++)    
    { 
      WriteSet(i,display[i]); 
    } 
	for(i =0 ;i <16 ; i++)      
   { 
      P0 = ReadSet(i);  
	  delaynms(200);  
	   delaynms(200); 
   }  
}
