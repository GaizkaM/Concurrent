with Ada.Text_IO; use Ada.Text_IO;
with Ada.Integer_Text_IO; use Ada.Integer_Text_IO;
with Ada.Task_Identification;
with Ada.Unchecked_Deallocation;
with Ada.Calendar;
with Ada.Strings.Fixed;
with Ada.Numerics;
with Random_Int_Generator;
with def_monitor;
use def_monitor;

procedure P2 is
   -----
   -- Tipus protegit per a la SC
   -----
   monitor : Maquina;

   -- Array de noms dels clients
   subtype Nom_Clients_Type is String (1 .. 10);
   NOM_CLIENTS : array (1 .. 10) of Nom_Clients_Type := ("Aina      ", "Albert    ", "Bel       ", "Bernat    ",
                                                         "Carla     ", "Dani      ", "Emma      ", "Ferran    ",
                                                         "Gina      ", "Hugo      ");

   -----
   -- Constants i variables
   -----
   NUM_CONSUMICIONS : constant Integer := 10;
   Num_Refresc_Maquina : Integer := 0;
   ConsumicionsClient : Integer := 0;

   ----- SIMULACIÓ EQUILIBRADA
   --Num_Clients : Integer := 4;
   --Num_Reposadors : Integer := 4;
   -----
   ----- SIMULACIÓ ALEATORIA
   Num_Clients : Integer := Random_Int_Generator.Get_Random_Integer;
   Num_Reposadors : Integer := Random_Int_Generator.Get_Random_Integer;
   -----

   -----
   -- Especificacio de la tasca Client
   -----
   task type Client_Task is
      entry Start(Nom : in String; Consumicions: in Integer);
   end Client_Task;

   -----
   -- Cos de la tasca Client
   -----
   task body Client_Task is
      My_Nom : String (1 .. 10);
      My_Consumicions : Integer;
   begin
      accept Start (Nom : in String; Consumicions : in Integer) do
         My_Nom := Nom;
         My_Consumicions := Consumicions;
      end Start;

      Put_Line (My_Nom & " diu: Hola, avui faré" & Integer'Image(My_Consumicions) & " consumicions");

      -- Bucle que es repiteix el nombre de consumicions que necessita el client
      for I in 1 .. My_Consumicions loop
         -- Simulació d'agafar el refresc
         delay 0.2;
         -- Bloqueig al monitor per part del client
         monitor.clientLock;
         -- Verificam si hi ha reposadors
         if Num_Reposadors = 0 then
            Put_Line (My_Nom & " diu: No hi ha reposadors a la màquina, me'n vaig");
            monitor.clientUnlock;
            exit;
         end if;

         -- Verificam si hi ha refrescs a la màquina
         if Num_Refresc_Maquina > 0 then
            Num_Refresc_Maquina := Num_Refresc_Maquina - 1;
            Put_Line ("---------- " & My_Nom & " agafa el refresc " & Integer'Image (I) & "/" & Integer'Image (My_Consumicions) & " a la màquina en queden " & Integer'Image (Num_Refresc_Maquina));
         end if;
         -- Desbloqueig monitor
         monitor.clientUnlock;
      end loop;
      -- El client ha finalitzat
      Num_Clients := Num_Clients - 1;
      Put_Line (My_Nom & " acaba i se'n va, queden " & Integer'Image (Num_Clients) & " clients >>>>>>>>>");
   end Client_Task;



  -----
  -- Especificacio de la tasca Reposador
  -----
   task type Reposador_Task is
      entry Start(Identificador: Integer);
   end Reposador_Task;

   -----
   -- Cos de la tasca Reposador
   -----
   task body Reposador_Task is
      My_id : Integer;
      Refresc_Reposats: Integer :=0;
   begin
      accept Start (Identificador: Integer) do
         My_id := Identificador;
      end Start;

      Put_Line ("     El reposador " & Integer'Image (My_id) & " comença a treballar");

      -- Bucle mentres hi quedin clients
      while Num_Clients /= 0 loop
         -- Simulació de reposar
         delay 0.1;
         -- Bloqueig del monitor
         monitor.reposadorLock;
         -- Verifica si la màquina està plena
         if Num_Refresc_Maquina /= NUM_CONSUMICIONS then
            Refresc_Reposats:= NUM_CONSUMICIONS - Num_Refresc_Maquina;
            Num_Refresc_Maquina := NUM_CONSUMICIONS;
            Put_Line ("++++++++++ El reposador " & Integer'Image (My_id) & " reposa " & Integer'Image (Refresc_Reposats) & " refrescs, ara n'hi ha " & Integer'Image (NUM_CONSUMICIONS));
         end if;
         -- Desbloqueig monitor
         monitor.reposadorUnlock;
      end loop;

      -- El reposador ha finalitzat
      Put_Line ("++++++++++ El reposador" & Integer'Image (My_id) & " diu: No hi ha clients me'n vaig");
      Put_Line ("     El reposador " & Integer'Image (My_id) & " acaba i se'n va >>>>>>>>>");
   end Reposador_Task;

begin
   -- Inici simulació
   Put_Line ("Simulació amb " & Integer'Image (Num_Clients) & " clients i " & Integer'Image (Num_Reposadors) & " Reposadors");

   -----
   -- Arrays de les tasques
   -----
   declare
      Clients   : array (1 .. Num_Clients) of Client_Task;
      Reposadors : array (1 .. Num_Reposadors) of Reposador_Task;
   begin
      -- Inicialització dels fils clients
      for I in Clients'Range loop
         -- Càlcul random de les consumicions de cada client
         ConsumicionsClient := Random_Int_Generator.Get_Random_Integer;
         Clients(I).Start(NOM_CLIENTS(I),ConsumicionsClient);
      end loop;

      -- Inicialització dels fils reposadors
      for I in Reposadors'Range loop
         Reposadors(I).Start(I);
      end loop;
   end;
end P2;




